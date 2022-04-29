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
