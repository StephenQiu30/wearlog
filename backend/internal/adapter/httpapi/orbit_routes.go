package httpapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerOrbitOperations(api huma.API, handler *OrbitHandler) {
	huma.Register(api, authenticatedOperation(huma.Operation{
		OperationID:   "analyzeOrbitOutfit",
		Method:        http.MethodPost,
		Path:          "/orbit/analyze",
		Tags:          []string{"Orbit"},
		Summary:       "分析全身穿搭并解构单品与季节打标",
		DefaultStatus: http.StatusOK,
		MaxBodyBytes:  64 * 1024,
		Errors:        []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusServiceUnavailable, http.StatusInternalServerError},
	}), handler.analyze)

	huma.Register(api, authenticatedOperation(huma.Operation{
		OperationID:   "tryOnOrbitLook",
		Method:        http.MethodPost,
		Path:          "/orbit/looks/{look_id}/try-on",
		Tags:          []string{"Orbit"},
		Summary:       "保持原生姿态的 2D 虚拟换装",
		DefaultStatus: http.StatusOK,
		MaxBodyBytes:  64 * 1024,
		Errors:        []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusServiceUnavailable, http.StatusInternalServerError},
	}), handler.tryOn)

	huma.Register(api, authenticatedOperation(huma.Operation{
		OperationID:   "generateOrbit360",
		Method:        http.MethodPost,
		Path:          "/orbit/looks/{look_id}/generate-360",
		Tags:          []string{"Orbit"},
		Summary:       "生成 16 帧 360° 交互式关键帧精灵图",
		DefaultStatus: http.StatusOK,
		Errors:        []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusServiceUnavailable, http.StatusInternalServerError},
	}), handler.generate360)

	huma.Register(api, authenticatedOperation(huma.Operation{
		OperationID:   "getOrbitLook",
		Method:        http.MethodGet,
		Path:          "/orbit/looks/{look_id}",
		Tags:          []string{"Orbit"},
		Summary:       "获取已完成的穿搭立体资产与单品列表",
		DefaultStatus: http.StatusOK,
		Errors:        []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusServiceUnavailable, http.StatusInternalServerError},
	}), handler.getLook)
}
