# Alamo

### 介绍
按照配置文件执行作业，可以设置作业的：
- 执行命令
- 所需要的环境变量
- flags
- 参数
- 前置作业，本作业需要依赖的其他作业
- 后置作业，本作业执行后的事后处理

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
Alamo需要一个作业配置文件，可以用--config来指定，否则，从当前路径（.)、home路径（$HOME/.alamo）和etc（/etc/alamo）下找alamo.yml文件。
配置文件的示例如下：
```yml
jobs:
  all:
    pre_jobs:
      - backup1
  backup1:
    command: restic
    env_vars:
      - RESTIC_PASSWORD: abc123
    flags:
      - --exclude: nwq*
      - --tag: job_backup1
      - -r: d:\tmp\srv\repo1
    arguments:
      - backup
      - d:\tmp\n.txt
      - d:\tmp\man
    pre_jobs:
      - depjob1
    post_jobs:
      - forget1
  depjob1:
    command: restic
    env_vars:
      - RESTIC_PASSWORD: abc123
    flags:
      - -r: d:\tmp\srv\repo1
    arguments:
      - snapshots
  forget1:
    command: restic
    env_vars:
      - RESTIC_PASSWORD: abc123
    flags:
      - -r: d:\tmp\srv\repo1
      - --keep-last: 5
      # - --tag: job_backup1
    arguments:
      - forget
```
