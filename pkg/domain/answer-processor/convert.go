package answerprocessor

func convertProcessAnswerParamsToGetWeightMatrixParams(
	params ProcessAnswerParams,
) GetWeightMatrixParams {
	return GetWeightMatrixParams{
		QuestionID: params.QuestionID,
	}
}

func convertProcessAnswerParamsToGetOtherPersonAnswersParams(
	params ProcessAnswerParams,
) GetOtherPersonAnswersParams {
	return GetOtherPersonAnswersParams{
		QuestionID: params.QuestionID,
	}
}

func convertPersonsWeightsToUpdatePersonsWeightsParams(
	params ProcessAnswerParams,
	personWeights []RelatedPersonWeight,
) UpdatePersonsWeightsParams {
	return UpdatePersonsWeightsParams{
		PersonID:       params.PersonID,
		QuestionID:     params.QuestionID,
		RelatedWeights: personWeights,
	}
}
