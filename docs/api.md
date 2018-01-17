# API

Draft

## Master

### Nodes

#### Get nodes

`GET /nodes/`

### Create new node

`POST /nodes/`

```
params:
{
    "name": "node",
    "ip": "127.0.0.1"
}
```

### Jobs

#### Get jobs

`GET /jobs/`

#### Create new job

`POST /jobs/`

```
params:
{
    "name": "job1",
    "schedule": "*/25 * * * *",
    "tasks": [
        {
            "name": "task1",
            "command": "echo 'task1 done'",
            "next_tasks": "task2, task3"
        },
        {
            "name": "task2",
            "command": "echo 'task2 done'"
        },
        {
            "name": "task3",
            "command": "echo 'task3 done'"
        },
        {
            "name": "task4",
            "command": "echo 'task4 done'"
        }
    ]
}
```

#### Get the job detail

`GET /jobs/:id`


#### Update the job

`PUT /jobs/:id`

```
params:
{
    "name": "job1_updated",
    "schedule": "*/20 * * * *",
    "tasks": [
        {
            "id": "1",
            "name": "task1",
            "command": "echo 'task1 done'",
            "next_tasks": "task3_updated"
        },
        {
            "id": "2",
            "name": "task2",
            "command": "echo 'task3 done'"
        },
        {
            "id": "3",
            "name": "task3_updated",
            "command": "echo 'task3 done'"
        },
        {
            "name": "task5_new",
            "command": "echo 'task5 done'"
        }
    ]
}
```

#### Delete the job

`DELETE /jobs/:id`
