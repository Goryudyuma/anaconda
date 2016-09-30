package anaconda

import (
	"net/url"
	"strconv"
)

// CreateList implements /lists/create.json
func (a TwitterApi) CreateList(name, description string, v url.Values) (list List, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("name", name)
	v.Set("description", description)

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/create.json", v, &list, _POST, response_ch}
	return list, (<-response_ch).err
}

// AddUserToList implements /lists/members/create.json
func (a TwitterApi) AddUserToList(screenName string, listID int64, v url.Values) (users []User, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", screenName)

	var addUserToListResponse AddUserToListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create.json", v, &addUserToListResponse, _POST, response_ch}
	return addUserToListResponse.Users, (<-response_ch).err
}

// GetListsOwnedBy implements /lists/ownerships.json
// screen_name, count, and cursor are all optional values
func (a TwitterApi) GetListsOwnedBy(userID int64, v url.Values) (lists []List, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("user_id", strconv.FormatInt(userID, 10))

	var listResponse ListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/ownerships.json", v, &listResponse, _GET, response_ch}
	return listResponse.Lists, (<-response_ch).err
}

func (a TwitterApi) GetListTweets(listID int64, includeRTs bool, v url.Values) (tweets []Tweet, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/statuses.json", v, &tweets, _GET, response_ch}
	return tweets, (<-response_ch).err
}

// Implement /lists/members by list_id
func (a TwitterApi) GetListMembers(listID int64, v url.Values) (users UserCursor, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("list_id", strconv.FormatInt(listID, 10))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members.json", v, &users, _GET, response_ch}
	return users, (<-response_ch).err
}

// Implement /lists/members by list_slug with (owner_id OR owner_screen_name)
func (a TwitterApi) GetListMembersBySlug(listname string, owner_screen_name string, owner_id int64, v url.Values) (users UserCursor, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("slug", listname)
	if owner_screen_name != "" {
		v.Set("owner_screen_name", owner_screen_name)
	}
	if owner_id != 0 {
		v.Set("owner_id", strconv.FormatInt(owner_id, 10))
	}

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members.json", v, &users, _GET, response_ch}
	return users, (<-response_ch).err
}

// Implement /lists/members/destroy by list_slug with (owner_id OR owner_screen_name) and screen_name
func (a TwitterApi) RemoveMemberFromList(listname string, remove_user_screen_name string, remove_user_id int64, owner_screen_name string, owner_id int64, v url.Values) (users UserCursor, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("slug", listname)

	if remove_user_screen_name != "" {
		v.Set("screen_name", remove_user_screen_name)
	}
	if remove_user_id != 0 {
		v.Set("user_id", strconv.FormatInt(remove_user_id, 10))
	}
	if owner_screen_name != "" {
		v.Set("owner_screen_name", owner_screen_name)
	}
	if owner_id != 0 {
		v.Set("owner_id", strconv.FormatInt(owner_id, 10))
	}

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/destroy.json", v, &users, _POST, response_ch}
	return users, (<-response_ch).err
}

// AddUserToListIds implements /lists/members/create_all.json
func (a TwitterApi) AddUserToListIds(user_ids []int64, slug string, owner_id int64, v url.Values) (users []User, err error) {
	resusers := make([]User, 0, 0)
	if len(user_ids) > 100 {
		resusers, err := a.AddUserToListIds(user_ids[100:], slug, owner_id, v)
		if err != nil {
			return resusers, err
		}
	}
	length := 100
	if len(user_ids) < length {
		length = len(user_ids)
	}
	user_ids = user_ids[:length]
	if v == nil {
		v = url.Values{}
	}

	var user_ids_string string
	for i, v := range user_ids {
		user_ids_string += strconv.FormatInt(v, 10)
		if i != len(user_ids)-1 {
			user_ids_string += ","
		}
	}

	v.Set("slug", slug)
	v.Set("user_id", user_ids_string)
	v.Set("owner_id", strconv.FormatInt(owner_id, 10))

	var addUserToListResponse AddUserToListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create_all.json", v, &addUserToListResponse, _POST, response_ch}
	return append(addUserToListResponse.Users, resusers...), (<-response_ch).err
}

// RemoveUserToListIds implements /lists/members/destroy_all.json
func (a TwitterApi) RemoveUserToListIds(user_ids []int64, slug string, owner_id int64, v url.Values) (users []User, err error) {
	resusers := make([]User, 0, 0)
	if len(user_ids) > 100 {
		resusers, err := a.RemoveUserToListIds(user_ids[100:], slug, owner_id, v)
		if err != nil {
			return resusers, err
		}
	}
	length := 100
	if len(user_ids) < length {
		length = len(user_ids)
	}
	user_ids = user_ids[:length]
	if v == nil {
		v = url.Values{}
	}

	var user_ids_string string
	for i, v := range user_ids {
		user_ids_string += strconv.FormatInt(v, 10)
		if i != len(user_ids)-1 {
			user_ids_string += ","
		}
	}

	v.Set("slug", slug)
	v.Set("user_id", user_ids_string)
	v.Set("owner_id", strconv.FormatInt(owner_id, 10))

	var addUserToListResponse AddUserToListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/destroy_all.json", v, &addUserToListResponse, _POST, response_ch}
	return append(addUserToListResponse.Users, resusers...), (<-response_ch).err
}
