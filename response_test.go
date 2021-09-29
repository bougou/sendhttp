package sendhttp

type fakeResponse struct {
	*BaseResponse `json:"-"`
	MSG           string `json:"msg"`
}

func newFakeResponse() *fakeResponse {
	return &fakeResponse{
		BaseResponse: NewBaseResponse(),
	}
}
