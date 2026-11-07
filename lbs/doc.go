// Package lbs provides a client for Tencent Location Service geocoding APIs.
//
// Design constraints:
//   - AppKey and timeout come from caller-injected Config; no viper in the library.
//   - NewClient returns an error when AppKey is missing.
//   - IO methods accept context.Context.
package lbs
