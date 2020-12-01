dir /data/sentinel
sentinel monitor mymaster {{ .Env.REDIS_MASTER_IP }} 6379 2
sentinel down-after-milliseconds mymaster 30000
sentinel parallel-syncs mymaster 1
sentinel auth-pass mymaster test@hago
sentinel failover-timeout mymaster 180000
sentinel deny-scripts-reconfig yes
port 26379
logfile "/data/sentinel/log"
