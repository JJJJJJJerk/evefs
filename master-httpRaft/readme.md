

### Http Server With Raft

- RESTful APIs
    - JWT Authorization
    - create bucket
    - remove bucket
    - file write
        - permission
        - create snow flake id
    - file read
        - read by file_name through mysql with redis cache
    - file resize
    - file delete
    - file download through URL
- store needles in mysql
- snow flake file ID
- Admin Backend UI
    - health 状态
    - prometheus 监控
    - 触发空间清理
    - store status
    - hay-volume status
- http Proxy to Raft Leader
- raft 分布式系统(optional/not urgent)

### required software

- mysql 
- redis 