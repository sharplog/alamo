# Alamo

### Intro
A tool to execute job according to configuration. Config items include:
- Command
- Arguments
- Flags
- Environment variables
- Work directory
- Prepositive jobs(dependencies)
- Post-processing jobs
- Post-failing jobs
- Redirection of stdin, stdout and stderr
- Dependent platform, starting with ^ means not matched

### Usage
#### Execute jobs
Specify one or more job names which are defined in configuration file when run alamo。
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
#### Configuration
Alamo needs a configuration file which can be specified by *--config*. By default, alamo looks for alamo.yml from current directory(.), home path(~/.), and etc(/etc/alamo).

One sample of configuration：
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
