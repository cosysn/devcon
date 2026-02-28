# E2E 测试用例详细文档

## TC_E2E001: full-flow-registry-feature-build

- **用例名称**: 完整流程：镜像仓库 + feature + 构建
- **测试文件**: tests/cases/e2e_full_flow.sh
- **测试内容**: 端到端测试完整流程

### 测试步骤

1. **构建 devcon CLI**
   ```bash
   go build -o devcon ./cmd/devcon
   ```

2. **启动本地镜像仓库**
   ```bash
   docker run -d --name devcon-test-registry -p 5000:5000 --restart=always registry:2
   ```

3. **创建测试 feature**
   ```bash
   # 创建目录和文件
   mkdir -p /tmp/devcon-e2e-test/feature

   # 创建 devcontainer-feature.json
   cat > /tmp/devcon-e2e-test/feature/devcontainer-feature.json << 'EOF'
   {
       "id": "test",
       "name": "Test Feature",
       "version": "1.0.0"
   }
   EOF

   # 创建 install.sh
   cat > /tmp/devcon-e2e-test/feature/install.sh << 'EOF'
   #!/bin/bash
   echo "Test feature installed!"
   EOF
   chmod +x /tmp/devcon-e2e-test/feature/install.sh
   ```

4. **发布 feature 到本地仓库**
   ```bash
   ./devcon features publish /tmp/devcon-e2e-test/feature --reg localhost:5000/test-feature:latest
   ```

5. **验证 feature 已发布**
   ```bash
   curl -s http://localhost:5000/v2/_catalog
   # 期望输出: {"repositories":["test-feature"]}
   ```

6. **创建 devcontainer 配置并构建**
   ```bash
   # 创建 devcontainer.json
   mkdir -p /tmp/devcon-e2e-test/project/.devcontainer
   cat > /tmp/devcon-e2e-test/project/.devcontainer/devcontainer.json << 'EOF'
   {
       "image": "alpine:latest",
       "features": {
           "localhost:5000/test-feature:latest": {}
       }
   }
   EOF

   # 构建
   ./devcon build /tmp/devcon-e2e-test/project
   # 期望输出: Image built: alpine:latest
   ```

7. **清理测试资源**
   ```bash
   # 清理容器
   docker rm -f devcon-test-container 2>/dev/null || true

   # 清理镜像
   docker rmi localhost:5000/test-feature:latest 2>/dev/null || true
   docker rmi localhost:5000/devcon-test:latest 2>/dev/null || true

   # 清理测试目录
   rm -rf /tmp/devcon-e2e-test
   ```

### 期望结果

- 所有步骤执行成功
- 镜像仓库成功启动
- Feature 成功发布到本地仓库
- 构建命令成功执行
- 清理后无残留容器和镜像

### 验证命令

```bash
make test-e2e-full-flow
```

或直接运行脚本：

```bash
./tests/cases/e2e_full_flow.sh
```
