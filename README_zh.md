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
jobs:
  all:
    pre_jobs:
      - backup1
  backup1:
    command: restic
    env_vars: &pwd
      - RESTIC_PASSWORD: abc123
    flags:
      - --exclude: nwq*
      - --tag: job_backup1
      - -r: &repo d:\tmp\srv\repo1
    arguments:
      - backup
      - d:\tmp\n.txt
      - d:\tmp\man
    pre_jobs:
      - depjob1
    post_jobs:
      - forget1
    fail_jobs:
    work_dir: d:\tmp
    stdin: in.txt
    stdout: out.txt
    stderr: err.txt
  depjob1:
    command: restic
    env_vars: *pwd
    flags:
      - -r: *repo
    arguments:
      - snapshots
  forget1:
    command: restic
    env_vars: *pwd
    flags:
      - -r: *repo
      - --keep-last: 5
      # - --tag: job_backup1
    arguments:
      - forget
```
