from fastapi import Request, status
from fastapi.responses import JSONResponse

from app.core.exceptions import FetchScheduleException, ParsingException


async def fetch_schedule_exception_handler(
    _request: Request, exc: FetchScheduleException
):
    return JSONResponse(
        status_code=status.HTTP_502_BAD_GATEWAY, content={"detail": str(exc)}
    )


async def parsing_exception_handler(_request: Request, exc: ParsingException):
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": str(exc)},
    )


async def url_not_found_exception(_request: Request, exc: ParsingException):
    return JSONResponse(
        status_code=status.HTTP_404_NOT_FOUND,
        content={"detail": str(exc)},
    )
