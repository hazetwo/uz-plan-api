from fastapi import Request, status
from fastapi.responses import JSONResponse

from app.api.core.exceptions import ParsingException, ScheduleFetchException


async def schedule_fetch_exception_handler(
    _request: Request, exc: ScheduleFetchException
):
    return JSONResponse(
        status_code=status.HTTP_502_BAD_GATEWAY, content={"detail": str(exc)}
    )


async def parsing_exception_handler(_request: Request, exc: ParsingException):
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": str(exc)},
    )
