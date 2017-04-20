## Strategy Document
Metric = (idle instance)/(total instance)
```
{
    [
        {
           "Metric":0.8,
           "Cmd":"scale -e constraint:node==IatScale -e affinity:GPU!=true -e TASK_TYPE=start -f app==iat -n 5"
        },
        {
           "Metric":0.2,
           "Cmd":"scale -e constraint:node==IatScale -e TASK_TYPE=stop -f app==iat -n -5"
        },
    ]
}
```
