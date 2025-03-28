# GitHub Actions 工作流

本目录包含项目的 GitHub Actions 工作流配置文件，用于自动化测试和部署流程。

## 可用工作流

### Go Tests (`go-test.yml`)

这个工作流在每次推送到 main/master 分支或提交 PR 时运行。它执行以下任务:

#### 测试和代码质量检查 Job

1. **代码格式检查**：使用 `gofmt` 检查代码格式是否规范
2. **静态分析**：使用 `go vet` 进行静态代码分析
3. **标准测试**：运行所有包的单元测试
4. **竞态检测**：使用 `-race` 标志运行测试以检测潜在的数据竞争问题
5. **代码覆盖率分析**：生成测试覆盖率报告
6. **上传覆盖率报告**：将覆盖率报告上传到 Codecov 服务

#### 构建 Job

1. **代码构建**：编译所有包确保代码可以成功构建
2. **示例构建**：编译 `examples` 目录下的所有示例程序

## 工作流状态

查看项目根目录的 README.md 文件中的状态徽章，可以了解最新的构建状态。

## 本地运行测试

以下命令与 CI 中运行的命令相同，可以在本地执行以验证代码:

```bash
# 检查代码格式
gofmt -l .

# 静态分析
go vet ./...

# 运行测试
go test -v ./...

# 运行带竞态检测的测试
go test -race -v ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
``` 