package product

import (
	"context"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/errgroup"
)

// WarmupStage 预热阶段接口。后续新增模块实现此接口即可加入预热管线。
// 每个阶段独立负责自身的数据加载和写入，不共享上下文状态。
type WarmupStage interface {
	Name() string
	Warmup(ctx context.Context) (int, error)
}

// WarmupPipeline 串联多个预热阶段，各阶段并行执行。
//
// 使用方式：
//
//	pipe := product.NewWarmupPipeline()
//	pipe.Add(&BrandWarmupStage{svc: brandSvc})
//	pipe.Add(&CategoryWarmupStage{svc: catSvc})
//	pipe.Run(ctx)
type WarmupPipeline struct {
	stages []WarmupStage
}

func NewWarmupPipeline(stages ...WarmupStage) *WarmupPipeline {
	return &WarmupPipeline{stages: stages}
}

func (p *WarmupPipeline) Add(stage WarmupStage) {
	p.stages = append(p.stages, stage)
}

// Run 并行执行所有已注册的预热阶段
func (p *WarmupPipeline) Run(ctx context.Context) map[string]int {
	results := make(map[string]int, len(p.stages))
	g, egCtx := errgroup.WithContext(ctx)
	for _, stage := range p.stages {
		stage := stage
		g.Go(func() error {
			n, err := stage.Warmup(egCtx)
			if err == nil {
				results[stage.Name()] = n
			}
			return nil
		})
	}
	_ = g.Wait()
	return results
}

// ── SPU 预热 ────────────────────────────────────────

// WarmupCache 全量预热 SPU 到 Bloom Filter + L2 + L1
// 数据只从 DB 加载一次，三个目标并行写入
func (s *SpuService) WarmupCache(ctx context.Context) (int, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}
	if len(all) == 0 {
		return 0, nil
	}

	g, _ := errgroup.WithContext(ctx)

	// Stage 1: Bloom Filter
	g.Go(func() error {
		ids := make([]int64, len(all))
		for i := range all {
			ids[i] = all[i].ID
		}
		s.bloomFilter.addAll(ids)
		return nil
	})

	// Stage 2: Redis Pipeline
	if s.rdb != nil {
		g.Go(func() error {
			pipe := s.rdb.Pipeline()
			for i := range all {
				data, marshalErr := sonic.Marshal(&all[i])
				if marshalErr != nil {
					continue
				}
				pipe.Set(ctx, cacheKeySPU(all[i].ID), data, spuEntityTTL)
			}
			_, err := pipe.Exec(ctx)
			return err
		})
	}

	// Stage 3: L1 本地缓存
	g.Go(func() error {
		for i := range all {
			s.localCache.warmupSingle(all[i].ID, &all[i])
		}
		return nil
	})

	return len(all), g.Wait()
}
