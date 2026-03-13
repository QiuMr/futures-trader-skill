# 发布到 GitHub 的完整指南

## 第一步：创建 GitHub 仓库

1. 登录你的 GitHub 账号
2. 点击右上角的 `+` 按钮
3. 选择 `New repository`
4. 填写仓库信息：
   - Repository name: `futures-trader-skill` (或你喜欢的名称)
   - Description: `Gate.io 期货交易 CLI 工具技能`
   - Public (必须选择公开)
   - 勾选 "Initialize this repository with a README"
5. 点击 `Create repository`

## 第二步：上传文件

### 方法一：使用 GitHub Web 界面（简单）

1. 进入你的仓库页面
2. 点击 "Add file" → "Upload files"
3. 拖拽或选择以下文件：
   - `SKILL.md`
   - `README.md`
   - `futures-trader.txt`
   - `futures-trader-linux-amd64.txt`
4. 点击 "Commit changes"

### 方法二：使用 Git 命令（推荐）

```bash
# 克隆你的仓库
git clone https://github.com/[你的GitHub用户名]/[仓库名].git
cd [仓库名]

# 复制文件
copy ..\go技能测试\SKILL.md .
copy ..\go技能测试\README.md .
copy ..\go技能测试\futures-trader.txt .
copy ..\go技能测试\futures-trader-linux-amd64.txt .

# 添加并提交
git add .
git commit -m "Initial commit: Gate.io futures trader skill"

# 推送到 GitHub
git push origin main
```

## 第三步：验证发布

1. 访问你的仓库页面
2. 确保文件都已上传
3. 确保仓库是 Public（公开）状态

## 第四步：分享你的技能

用户可以通过以下命令安装你的技能：

```bash
npx skills add [你的GitHub用户名]/[仓库名]
```

或者指定子目录：

```bash
npx skills add [你的GitHub用户名]/[仓库名]@main
```

## 文件结构

你的 GitHub 仓库应该包含：

```
[仓库名]/
├── SKILL.md          # 技能配置文件（必须）
├── README.md         # 说明文档（推荐）
├── futures-trader.txt           # Windows 可执行文件
├── futures-trader-linux-amd64.txt # Linux 可执行文件
└── github-guide.md   # 本指南（可选）
```

## 注意事项

1. **必须保持仓库公开**，否则 skills.sh 无法访问
2. **SKILL.md 必须在根目录**，或者在安装时指定路径
3. **文件大小限制**：GitHub 单文件最大 100MB，你当前的文件约 2.7MB，完全符合要求
4. **版本更新**：修改后直接 `git push` 即可更新

## 常见问题

**Q: 仓库创建后可以改名吗？**  
A: 可以，但会影响用户的安装命令，建议起好名字

**Q: 可以添加更多技能吗？**  
A: 可以，在同一个仓库里创建不同子目录即可

**Q: 如何更新技能？**  
A: 修改文件后 `git add . && git commit -m "更新说明" && git push` 即可
