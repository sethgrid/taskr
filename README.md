## Taskr

Maybe a better name for this is jrnl. Log a journal entry of tasks that you have completed with labels from the command line.

```
$ taskr -h
Usage of taskr:
  -db_name=".taskr.db": full path. defaults to .taskr.db
  -l="default": set a list of labels allowing you to group your message, comma delimited
  -label="": create a new label to use for messages
  -m="": a message entry to log
  -show=false: set to true to show all entries
  -show_labels=false: set to true to show all registered labels
```

Example entry:

```
$ taskr -label project_x
$ taskr -label project_y
$ taskr -label work
$ taskr -label personal

$ taskr -l project_x,work -m "contacted so-and-so and and changed project requirements to include authentication"

$ taskr -show
MON/DAY/YEAR - contacted so-and-so and and changed project requirements to include authentication [project_x,work]
```