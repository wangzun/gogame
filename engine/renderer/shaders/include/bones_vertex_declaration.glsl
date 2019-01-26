#ifdef BONE_INFLUENCERS
    #if BONE_INFLUENCERS > 0
	uniform mat4 mBones[TOTAL_BONES];
    attribute vec4 matricesIndices;
    attribute vec4 matricesWeights;
//    #if BONE_INFLUENCERS > 4
//        in vec4 matricesIndicesExtra;
//        in vec4 matricesWeightsExtra;
//    #endif
    #endif
#endif
