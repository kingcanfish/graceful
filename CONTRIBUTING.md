# Contributing to Graceful

感谢您对 Graceful 项目感兴趣！我们欢迎任何形式的贡献，包括但不限于：

- 报告问题和建议
- 提交代码改进
- 改进文档
- 添加测试用例

## 开发流程

1. Fork 项目到您的 GitHub 账号
2. 克隆项目到本地：`git clone https://github.com/YOUR_USERNAME/graceful.git`
3. 创建新的分支：`git checkout -b feature/your-feature-name`
4. 进行代码修改
5. 运行测试：`go test -v -race ./...`
6. 提交代码：`git commit -m "feat: your commit message"`
7. 推送到远程：`git push origin feature/your-feature-name`
8. 创建 Pull Request

## 代码规范

- 遵循 Go 标准代码规范
- 所有新代码必须包含适当的测试
- 提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范
- 保持代码简洁，避免不必要的复杂性
- 确保所有测试通过，且代码覆盖率不降低

## Pull Request 规范

1. PR 标题应简洁明了地描述改动
2. 在 PR 描述中详细说明改动内容和原因
3. 确保 CI 检查通过
4. 如果修复了 issue，请在 PR 描述中引用相关 issue

## 开发建议

1. 在开始大型改动前，建议先创建 issue 讨论
2. 保持改动范围适中，便于审查
3. 及时同步上游仓库的更新
4. 多写注释，特别是对复杂逻辑的解释

## 许可证

通过提交 PR，您同意将您的代码按照项目的 MIT 许可证开源。

## 获取帮助

如果您在开发过程中遇到任何问题，欢迎创建 issue 咨询。