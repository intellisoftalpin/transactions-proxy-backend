package consts

const SessionLifetimeInHours = 24

const SessionsEmptifyPeriodInHours = 24

const CSuccessStatus = "Success"
const CErrorStatus = "Error"

const MessageOtherError = "Other."

const MessageNewUserAndNewSessionError = "Error while creating new session for new user."
const MessageNewUserAndNewSessionSuccess = "New user successful added. New session created."

const MessageUserExists = "User exists. "

const MessageSessionExpiredNewSessionError = "Old session expired. Error while creating new session."
const MessageSessionExpiredNewSessionSuccess = "Old session expired. New session created."

const MessageSessionNotExpiredOldKeyError = "Session not expired. Error while getting authorization key."
const MessageSessionNotExpiredOldKeySuccess = "Session not expired. Return authorization key."

const MessageSessionNotExistsNewSessionError = "Session not exist. Error while creating new session."
const MessageSessionNotExistsNewSessionSuccess = "Session not exist. New session created."

const MessageSessionNotExistsOrExpiredNewSessionError = "Session not exist or expired. Error while creating new session."
const MessageSessionNotExistsOrExpiredNewSessionSuccess = "Session not exist or expired. New session created."

const MessageDeleteTransactionError = "Error while deleting user`s transaction."
const MessageDeleteTransactionSuccess = "User`s transaction deleted successful."

const CTransactionStatusDraft = "draft"
const CTransactionStatusPrepared = "prepared"
const CTransactionStatusSubmitted = "submitted"
const CTransactionStatusSuccess = "success"
const CTransactionStatusFailed = "failed"

const (
	CErrorInvalidSessionID     = "invalidSessionID"
	CErrorsInternalError       = "internalError"
	CErrorInvalidTransactionID = "invalidTransactionID"
)

const (
	CActiveTransactionStepSubmit           = "submit"           // 1
	CActiveTransactionStepUpdateStatus     = "updateStatus"     // 2
	CActiveTransactionStepSend             = "send"             // 3
	CActiveTransactionStepOnlyUpdateStatus = "onlyUpdateStatus" // 4
)

const CSyncProgressStatusReady = "ready"
const CErrorWalletNetworkIsNotReady = "walletNetworkIsNotReady"
