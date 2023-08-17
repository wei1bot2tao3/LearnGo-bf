--1.检查是不是你的锁
--2.删除
--KEYS[1]是你的分布式锁的key
--ARGV[1]是你预期在redis里面的value
if redis.call('get',KEYS[1])==ARGV[1] then
    --  确实是你的锁
    return redis.call('del',KEYS[1])
else
    --不是你的锁
    return 0

end