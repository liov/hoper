---
title: 复合主键
layout: page
---

将多个字段设置为 primary key 以启用复合主键

```go
type Product struct {
   ID           string `gorm:"primaryKey"`
   LanguageCode string `gorm:"primaryKey"`
   Code         string
   Name         string
}
```