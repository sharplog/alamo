# Alamo

### 介绍
按照配置文件执行作业，可以设置作业的：
- 执行命令
- 参数
- flags
- 所需要的环境变量
- 工作目录
- 前置作业，本作业需要依赖的其他作业
- 后置作业，本作业执行后的事后处理
- 失败后作业，在本作业执行失败后执行
- 标准输入、标准输出和错误输出的重定向
- 运行平台, 以^开关意味不匹配

### 使用说明
#### 运行作业
运行时指定一个或多个作业名称，作业是在配置文件中配置的。
```bash
$ alamo -h
See https://gitee.com/logsharp/alamo for documentation.

Usage:
  alamo [flags] job [job] ...

Flags:
  -c, --config string      config file
  -h, --help               help for alamo
  -l, --log-level string   log level (fatal|error|warn|info|trace) (default "info")
```
#### 配置文件
Alamo需要一个作业配置文件，可以用--config来指定，否则，从当前路径（.)、home路径（~/.alamo）和etc（/etc/alamo）下找alamo.yml文件。

配置文件的示例如下：
```yml
# alamo config example

jobs:
  backup:
    command: restic
    env_vars: &pwd
      - RESTIC_PASSWORD: abc123
    flags:
      - --exclude: nwq*
      - --tag: job_backup
      - -r: &repo d:\tmp\srv\repo1
    arguments:
      - backup
      - d:\tmp\alamo.gz
    pre_jobs:
      - compress
    post_jobs:
      - forget
      - rm1
      - rm2
    fail_jobs:
    work_dir: d:\tmp
    stdin: in.txt
    stdout: out.txt
    stderr: err.txt
  backup:
    command: restic
    env_vars: &pwd
      - RESTIC_PASSWORD: abc123
    flags:
      - --exclude: nwq*
      - --tag: job_backup
      - -r: &repo d:\tmp\srv\repo1
    arguments:
      - backup
      - d:\tmp\alamo.gz
    pre_jobs:
      - compress
    post_jobs:
      - forget
      - rm_w
      - rm_x
    fail_jobs:
    work_dir: d:\tmp
    stdin: in.txt
    stdout: out.txt
    stderr: err.txt
  compress:
    command: 7z
    arguments:
      - a
      - alamo
      - alamo.txt
      - -tgzip
    work_dir: d:\tmp
  forget:
    command: restic
    env_vars: *pwd
    flags:
      - -r: *repo
      - --keep-last: 5
      - --tag: job_backup
    arguments:
      - forget
  rm_w:
    command: cmd
    platform: windows   # regexp matched platform
    arguments: 
      - /C
      - del alamo.gz
    work_dir: d:\tmp
  rm_x:
    command: rm
    platform: ^windows   # regexp not matched platform
    arguments: alamo.gz
    work_dir: /tmp
```
