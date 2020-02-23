# Git 原理操作

修改一个文件并添加

```shell script
echo "version 3" | git hash-object -w --stdin # 写入更改内容
# 7170a5278f42ea12d4b6de8ed1305af8c393e756
git update-index --add --cacheinfo 100644 7170a5278f42ea12d4b6de8ed1305af8c393e756 test.txt \
    # 更改内容与文件相关联，存储进暂存区
git write-tree # 记录暂存区，创建内容树
5bee27fca16056be0a14169ffef5b3eceb360028
echo "second commit" | git commit-tree 5bee27fca16056be0a14169ffef5b3eceb360028 -p 59536ed44ea9400a12e787a0f2049f3cc9f1f955 \
    # 提交更改，指定内容树和base提交
b41872763ca59a88411509d4aabe4d4fd56ac78e
```

- 初加入
    - 对每个资源池内文件撕裂，分发到其他参与者
    - 每个有同一个资源的人有一样的该资源的checksum

- 常态
    - 监听资源池内文件变化
    - 发布增量更新内容，同时发布新checksum
