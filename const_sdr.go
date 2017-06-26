/*
Copyright (c) 2014 EOITek, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipmi

// section 33.9
const (
	SDR_OP_SUP_ALLOC_INFO   = (1 << 0)
	SDR_OP_SUP_RESERVE_REPO = (1 << 1)
	SDR_OP_SUP_PARTIAL_ADD  = (1 << 2)
	SDR_OP_SUP_DELETE       = (1 << 3)
	SDR_OP_SUP_NON_MODAL_UP = (1 << 5)
	SDR_OP_SUP_MODAL_UP     = (1 << 6)
	SDR_OP_SUP_OVERFLOW     = (1 << 7)
)
