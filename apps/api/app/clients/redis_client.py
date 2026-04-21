import redis.asyncio as redis

from app.config.settings import settings

RedisClient = redis.Redis


def create_redis_client() -> RedisClient:
    return redis.from_url(settings.redis.URL)
