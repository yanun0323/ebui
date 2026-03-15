package ebui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPreload(t *testing.T) {
	suite.Run(t, new(PreloadSuite))
}

type PreloadSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *PreloadSuite) SetupSuite() {
	su.ctx = context.Background()
}

func (su *PreloadSuite) TestPreload() {
	testCases := []struct {
		desc     string
		view     SomeView
		expected preloadData
	}{
		{
			desc:     "ZStack",
			view:     ZStack(),
			expected: preloadData{},
		},
		{
			desc: "ZStack with child",
			view: ZStack(Rectangle()),
			expected: preloadData{
				IsInfWidth:  true,
				IsInfHeight: true,
			},
		},
		{
			desc: "ZStack with framed child",
			view: ZStack(Rectangle().Frame(Bind(NewSize(100, 100)))),
			expected: preloadData{
				FrameSize: NewSize(100, 100),
			},
		},
		{
			desc: "ZStack with framed padding child",
			view: ZStack(Rectangle().Frame(Bind(NewSize(100, 100))).Padding(Bind(CGInset{10, 10, 10, 10}))),
			expected: preloadData{
				FrameSize: NewSize(120, 120),
			},
		},
		{
			desc: "ZStack with framed padding child and padding",
			view: ZStack(
				Rectangle().Frame(Bind(NewSize(100, 100))).Padding(Bind(CGInset{10, 10, 10, 10})),
			).Padding(Bind(CGInset{10, 10, 10, 10})),
			expected: preloadData{
				FrameSize: NewSize(120, 120),
				Padding:   NewInset(10, 10, 10, 10),
			},
		},
		{
			desc: "ZStack with framed padding child and padding and border",
			view: ZStack(
				Rectangle().Frame(Bind(NewSize(100, 100))).Padding(Bind(CGInset{10, 10, 10, 10})).Border(Bind(CGInset{10, 10, 10, 10})),
			).Padding(Bind(CGInset{10, 10, 10, 10})).Border(Bind(CGInset{10, 10, 10, 10})),
			expected: preloadData{
				FrameSize: NewSize(140, 140),
				Padding:   NewInset(10, 10, 10, 10),
				Border:    NewInset(10, 10, 10, 10),
			},
		},
	}
	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			actual, _ := tc.view.preload(nil)
			su.Equal(tc.expected.FrameSize, actual.FrameSize)
			su.Equal(tc.expected.IsInfWidth, actual.IsInfWidth)
			su.Equal(tc.expected.IsInfHeight, actual.IsInfHeight)
			su.Equal(tc.expected.Padding, actual.Padding)
			su.Equal(tc.expected.Border, actual.Border)
		})
	}
}
