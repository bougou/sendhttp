package sendhttp

type fakeRequest struct {
	*BaseRequest
}

func newFakeRequest() *fakeRequest {
	return &fakeRequest{
		BaseRequest: NewBaseRequest(),
	}

}
