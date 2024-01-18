# ehole指纹识别重构版

## 优化功能

1. 只支持url或文件识别
2. 改用ants管理goroutine
3. 改用req代替默认的http发包
4. 汇总去重ehole_magic和finger的指纹



# 使用

```bash
go build main.go

```

![img.png](img/img1.png)


![img.png](img/img2.png)