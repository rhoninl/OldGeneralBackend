FSM
```mermaid
flowchart TB
    running
    pending
    success
    failed
    resurrect

    pending --> running
    running --> resurrect
    running --> failed
    resurrect --> failed
    resurrect --> running
    resurrect --> success
    running --> success
    running --> running
```