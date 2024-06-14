/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ReferralCode struct {
	// Referral code itself, unique identifier
	Id string `json:"id"`
	// Status of the code, belonging to this user (referrer):   1. active: the code is not used yet by another user (referee)   2. banned: the referrer's country (known after scanning passport)      is not allowed to participate in the referral program   3. limited: the limit of reserved tokens in the referrer's country is reached   4. awaiting: the code is used by referee who has scanned passport, but the referrer hasn't yet   5. rewarded: the code is used, both referee and referrer have scanned passports   6. consumed: the code is used by referee who has not scanned passport yet  The list is sorted by priority. E.g. if the referee has scanned passport, but referrer's country has limit reached, the status would be `limited`.
	Status string `json:"status"`
}
