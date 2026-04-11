---
name: "code-standard"
description: "检查项目代码是否符合项目的代码规范和最佳实践。当用户提交代码前、进行代码审查时或项目重构时调用。"
---

# 代码规范检查

## 技能信息

- **技能名称**: 代码规范检查
- **技能ID**: code-standard
- **描述**: 检查项目代码是否符合项目的代码规范和最佳实践
- **版本**: 1.0.0
- **作者**: System

## 适用场景

- 代码提交前的规范检查
- 代码审查过程中的规范验证
- 项目重构时的规范统一

## 代码规范

### 1. 项目结构

- **cmd目录**: 存放各服务的入口文件
- **internal目录**: 按服务划分，包含服务内部代码
  - **api**: 包含HTTP处理、路由和DTO
  - **domain**: 包含领域模型和仓库
  - **service**: 包含业务逻辑
  - **app**: 包含应用初始化和启动逻辑
  - **mq**: 包含消息队列相关代码
- **pkg目录**: 存放公共代码，供多个服务使用
- **configs目录**: 存放配置文件
- **scripts目录**: 存放脚本文件

### 2. 命名规范

- **包名**: 使用小写字母，简短清晰
- **结构体名**: 使用驼峰命名法，首字母大写
- **方法名**: 使用驼峰命名法，首字母大写表示可导出
- **变量名**: 使用驼峰命名法，首字母小写
- **常量名**: 使用全大写，单词间用下划线分隔
- **JSON标签**: 使用蛇形命名法（小写字母，单词间用下划线分隔）
- **数据库表名**: 使用蛇形命名法，复数形式

### 3. 代码风格

- **缩进**: 使用4个空格进行缩进
- **换行**: 每行不超过120个字符
- **括号**: 左括号不单独占一行
- **注释**: 使用中文注释，清晰说明代码功能
- **空行**: 函数间使用空行分隔，逻辑块间使用空行分隔

### 4. 错误处理

- **错误返回**: 函数应返回(error)作为最后一个返回值
- **错误处理**: 使用c.Error(err)处理HTTP请求中的错误
- **错误传递**: 不吞掉错误，应向上传递或处理
- **错误日志**: 重要错误应记录日志

### 5. 依赖管理

- **依赖声明**: 使用go.mod管理依赖
- **依赖版本**: 固定依赖版本，避免使用latest
- **依赖导入**: 按标准库、第三方库、本地包的顺序导入

### 6. 数据库操作

- **ORM框架**: 使用gorm作为ORM框架
- **主键**: 使用UUID作为主键
- **时间字段**: 包含CreatedAt、UpdatedAt、DeletedAt字段
- **软删除**: 使用gorm.DeletedAt实现软删除
- **钩子函数**: 使用BeforeCreate等钩子函数生成UUID

### 7. API设计

- **路由结构**: 采用RESTful风格
- **API版本**: 在URL中包含版本号（如/v1/）
- **请求参数**: 使用DTO结构体绑定请求参数
- **响应格式**: 使用统一的响应格式（response.Success）
- **API文档**: 使用swaggo生成API文档

### 8. 性能优化

- **分页查询**: 支持分页和筛选
- **连接池**: 使用数据库连接池
- **缓存**: 合理使用缓存
- **并发**: 合理使用goroutine

### 9. 安全规范

- **密码存储**: 使用bcrypt加密存储密码
- **JWT验证**: 使用JWT进行身份验证
- **输入验证**: 对所有输入进行验证
- **SQL注入**: 使用参数化查询，避免SQL注入

### 10. 测试规范

- **单元测试**: 为关键函数编写单元测试
- **集成测试**: 为服务间交互编写集成测试
- **测试覆盖率**: 保持较高的测试覆盖率

## 检查规则

### 1. 项目结构检查

- [x] 检查项目是否遵循标准目录结构
- [x] 检查服务目录结构是否一致
- [x] 检查公共代码是否放在pkg目录

### 2. 命名规范检查

- [x] 检查包名是否使用小写字母
- [x] 检查结构体名是否使用驼峰命名法
- [x] 检查方法名是否使用驼峰命名法
- [x] 检查变量名是否使用驼峰命名法
- [x] 检查常量名是否使用全大写
- [x] 检查JSON标签是否使用蛇形命名法
- [x] 检查数据库表名是否使用蛇形命名法

### 3. 代码风格检查

- [x] 检查缩进是否使用4个空格
- [x] 检查每行是否不超过120个字符
- [x] 检查括号是否正确使用
- [x] 检查注释是否清晰完整
- [x] 检查空行是否合理使用

### 4. 错误处理检查

- [x] 检查函数是否返回error
- [x] 检查错误是否正确处理
- [x] 检查错误是否向上传递
- [x] 检查错误是否记录日志

### 5. 依赖管理检查

- [x] 检查go.mod是否存在
- [x] 检查依赖版本是否固定
- [x] 检查依赖导入顺序是否正确

### 6. 数据库操作检查

- [x] 检查是否使用gorm
- [x] 检查主键是否使用UUID
- [x] 检查是否包含时间字段
- [x] 检查是否使用软删除
- [x] 检查是否使用钩子函数

### 7. API设计检查

- [x] 检查路由是否遵循RESTful风格
- [x] 检查URL是否包含版本号
- [x] 检查是否使用DTO绑定参数
- [x] 检查是否使用统一响应格式
- [x] 检查是否使用swaggo生成文档

## 工具集成

### 1. 代码检查工具

- **golint**: 检查代码风格
- **go vet**: 检查代码潜在问题
- **staticcheck**: 静态代码分析

### 2. 代码格式化工具

- **gofmt**: 格式化代码
- **goimports**: 自动调整导入顺序

### 3. 测试工具

- **go test**: 运行测试
- **gocov**: 生成测试覆盖率报告

## 最佳实践

1. **代码复用**: 提取公共代码到pkg目录
2. **接口设计**: 使用接口定义服务间交互
3. **配置管理**: 使用viper管理配置
4. **日志管理**: 使用zap进行日志记录
5. **消息队列**: 使用RabbitMQ进行异步通信
6. **服务发现**: 使用gRPC进行服务间调用
7. **监控告警**: 实现系统监控和告警
8. **CI/CD**: 配置持续集成和持续部署

## 常见问题

### 1. 命名不一致

**问题**: 变量名、函数名、结构体名命名风格不一致

**解决方案**: 统一使用驼峰命名法，遵循项目规范

### 2. 错误处理不当

**问题**: 错误被吞掉或处理不当

**解决方案**: 正确返回和处理错误，使用c.Error(err)处理HTTP请求错误

### 3. 代码冗余

**问题**: 重复代码过多

**解决方案**: 提取公共代码到pkg目录，使用函数封装重复逻辑

### 4. 性能问题

**问题**: 数据库查询性能差，代码执行效率低

**解决方案**: 使用分页查询，合理使用缓存，优化SQL语句

### 5. 安全问题

**问题**: 密码存储不安全，存在SQL注入风险

**解决方案**: 使用bcrypt加密存储密码，使用参数化查询

## 总结

本代码规范技能旨在帮助开发团队保持代码的一致性和质量，提高代码的可维护性和可读性。通过遵循这些规范和最佳实践，开发团队可以更高效地协作，减少错误，提高代码质量。

## 示例代码

### 结构体定义示例

```go
type User struct {
	ID     string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Status int    `gorm:"type:tinyint;default:1" json:"status"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	UserInfo *UserInfo `gorm:"foreignKey:UserID" json:"user_info,omitempty"`
	Roles    []Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
```

### HTTP处理器示例

```go
// CreateInventory 创建库存
// @Summary 创建库存
// @Description 创建一个新的库存记录
// @Tags inventories
// @Accept json
// @Produce json
// @Param inventory body dto.CreateInventoryDTO true "库存信息"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	inventory, err := h.inventorySvc.CreateInventory(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishInventoryCreated(inventory)
	}
	response.Success(c, inventory)
}
```

### 服务层示例

```go
func (s *InventoryService) CreateInventory(ctx context.Context, req dto.CreateInventoryDTO) (*models.Inventory, error) {
	inventory := &models.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    1,
	}

	if err := s.repo.Create(ctx, inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}
```
