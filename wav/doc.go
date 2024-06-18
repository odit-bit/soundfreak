package wav

/*
	this package exist type thac could read and write wav file
	for remainder :
	maximum size of wav is 4gb:
	in raw audio sample could be represent by (8,16,32 bit) byte size
	wav has channel (pcm mode) data interleaved
	for every sample it represent by int32 (4 byte representation) or IEEE 754 (float64)

	for pcm example:
	sample1 24 bit = []byte{255,255,255}
	2 ch  = []byte{sample1-ch1, sample1-ch2, ......}

*/
