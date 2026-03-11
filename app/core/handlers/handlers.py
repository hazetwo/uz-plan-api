from fastapi import FastAPI, Request, status
from fastapi.responses import JSONResponse

from app.core.handlers.exceptions import (
    FetchScheduleException,
    GroupsDataException,
    ParsingException,
    UrlNotFoundException,
)


async def fetch_schedule_exception_handler(
    _request: Request, exc: Exception
) -> JSONResponse:
    assert isinstance(exc, FetchScheduleException)
    return JSONResponse(
        status_code=status.HTTP_502_BAD_GATEWAY,
        content={"detail": str(exc)},
    )


async def parsing_exception_handler(
    _request: Request, exc: Exception
) -> JSONResponse:
    assert isinstance(exc, ParsingException)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": str(exc)},
    )


async def url_not_found_exception_handler(
    _request: Request, exc: Exception
) -> JSONResponse:
    assert isinstance(exc, UrlNotFoundException)
    return JSONResponse(
        status_code=status.HTTP_404_NOT_FOUND,
        content={"detail": str(exc)},
    )


async def groups_data_exception_handler(
    _request: Request, exc: Exception
) -> JSONResponse:
    assert isinstance(exc, GroupsDataException)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": str(exc)},
    )


def register_exception_handlers(app: FastAPI) -> None:
    app.add_exception_handler(
        FetchScheduleException, fetch_schedule_exception_handler
    )
    app.add_exception_handler(ParsingException, parsing_exception_handler)
    app.add_exception_handler(
        UrlNotFoundException, url_not_found_exception_handler
    )
