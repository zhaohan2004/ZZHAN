---
colors:
  bg: "#f6fafb"
  bg-soft: "#edf6f8"
  text: "#203743"
  accent: "#206079"
  note-bg: "#e8f4f7"
  note-mark: "#f2d66b"
  dark-bg: "#14222b"
  dark-accent: "#a1d8eb"
omitted:
  typography: "System font stacks and roles are documented in Typography and owned by runtime CSS."
  rounded: "Existing component rules own geometry."
  spacing: "Existing responsive rules own layout."
  components: "The existing Vue components and shared stylesheet own implementation."
---

# CodeThink Blog Design System

## Overview

### Creative North Star

像在晴天的窗边翻看一本技术笔记。首页左侧用留白衬托站点名，右侧是一页带浅蓝页边线和柠檬黄书签的作者手记，下方进入文章索引。

### Product context and register

- 面向阅读中文技术文章的访客，支持浏览、筛选和长文阅读。
- 前台默认简体中文，文章与代码可中英混排。桌面和移动端均使用文档滚动。
- 适用于公开博客。主要视觉元素是取自后台终端配置的作者手记。
- 配色以静态色面、细边框和实色文字为主。首页避免发光装饰、漂浮动画和彩虹式色阶。
- 运行时以 `ZZHAN_web/user/src/styles/main.css` 的 `:root` 与 `html.dark` 为准。上方记录主色及深色对应值，Tailwind 通过变量读取同一来源，主题状态由现有 Pinia store 驱动。

## Colors

白色首页与雾蓝阅读底色拉开层次。深灰蓝正文配淡湖蓝色面，较深的湖蓝用于链接、焦点和主操作。深色使用夜蓝和浅湖蓝。实色按钮搭配各自主题的 `--on-accent`，成功、危险和警告由各自语义色表达。柠檬黄 `--note-mark` 用于手记和浏览器图标上的书签。

## Typography

`--font-display` 采用系统宋体与衬线字体，用于首页主标题、区块标题和文章标题。`--font-body` 使用系统无衬线字体，用于阅读、导航和表单。`--font-mono` 用于手记、日期与代码。字体使用本机回退栈，避免网络字体替换影响布局。

## Layout

共享容器上限 1200px。首页桌面为站点介绍与手记双栏，768px 以下按阅读顺序堆叠。文章封面保留媒体位，列表和正文继续使用文档滚动。窄屏的站点名称可缩短显示，首页标题仍完整展示。

## Elevation & Depth

用留白、浅蓝色面和细边线区分导航、正文与卡片。浮层和交互反馈可用低对比阴影。

## Shapes

文章卡片使用 14px 圆角，手记使用 16px 圆角。手记右上角有小书签，左侧页边线对齐文字。

## Components

### Foundational visual states

选中和聚焦使用主题强调色。键盘焦点使用独立轮廓，滚动条随主题切换并兼容强制配色。保留现有加载骨架及空状态。

### Buttons and actions

阅读主入口使用实色，辅助动作用描边。危险动作继续使用独立语义色。

### Navigation and data display

当前导航项使用下划线，窄屏菜单以色面标记当前项。分类名称保留可读文字，分类配置色只用于淡色图标底。未上传封面时，索引卡以分类和文章编号标识内容，使用淡湖蓝、雾紫或浅米色底。

### Forms and overlays

输入、弹层和提示沿用共享样式及其语义色。

### Iconography

沿用 Lucide 线条图标。阅读主入口同时保留文字。

### Motion

手记内容直接显示。动效使用已有的短反馈，并遵循系统的减少动态效果偏好。

### Content and data visualization

沿用后台配置的站点名称、简介、手记内容、文章封面与分类。日期、阅读量等元数据保留既有格式。

## Do's and Don'ts

- 用共享变量维护前台亮暗主题。
- 把柠檬黄留给小面积的书签。
- 避免装饰性渐变、循环装饰动画和低对比元数据。
