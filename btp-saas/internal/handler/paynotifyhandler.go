package handler

import (
	"net/http"

	"btp-saas/internal/logic"
	"btp-saas/internal/svc"
	"btp-saas/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayNotifyRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPayNotifyLogic(r.Context(), svcCtx)
		err := l.PayNotify(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Write([]byte("ok"))
			httpx.Ok(w)
		}
	}
}
