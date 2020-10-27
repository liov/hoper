package user

/*func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := NewMockUserServiceServer(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	r := rand.New(rand.NewSource(1))
	req := NewPopulatedActiveReq(r, false)
	m.EXPECT().
		Active(context.Background(), req).
		Return(response.NewPopulatedRawReply(r, false), nil).
		AnyTimes()
	res, _ := m.Active(context.Background(), req)
	log.Info(res)
}
*/
