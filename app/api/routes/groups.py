from fastapi import APIRouter, Request

router = APIRouter(prefix="/groups", tags=["groups"])


@router.get("/")
def get_groups(request: Request):
    return request.app.state.groups
