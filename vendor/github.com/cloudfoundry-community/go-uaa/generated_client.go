// Code generated by go-uaa/generator; DO NOT EDIT.

package uaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetClient with the given clientID.
func (a *API) GetClient(clientID string) (*Client, error) {
	u := urlWithPath(*a.TargetURL, fmt.Sprintf("%s/%s", ClientsEndpoint, clientID))
	client := &Client{}
	err := a.doJSON(http.MethodGet, &u, nil, client, true)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateClient creates the given client.
func (a *API) CreateClient(client Client) (*Client, error) {
	u := urlWithPath(*a.TargetURL, ClientsEndpoint)
	created := &Client{}
	j, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = a.doJSON(http.MethodPost, &u, bytes.NewBuffer([]byte(j)), created, true)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// UpdateClient updates the given client.
func (a *API) UpdateClient(client Client) (*Client, error) {
	u := urlWithPath(*a.TargetURL, ClientsEndpoint)
	created := &Client{}
	j, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = a.doJSON(http.MethodPut, &u, bytes.NewBuffer([]byte(j)), created, true)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// DeleteClient deletes the client with the given client ID.
func (a *API) DeleteClient(clientID string) (*Client, error) {
	if clientID == "" {
		return nil, errors.New("clientID cannot be blank")
	}
	u := urlWithPath(*a.TargetURL, fmt.Sprintf("%s/%s", ClientsEndpoint, clientID))
	deleted := &Client{}
	err := a.doJSON(http.MethodDelete, &u, nil, deleted, true)
	if err != nil {
		return nil, err
	}
	return deleted, nil
}

// ListClients with the given filter, sortBy, attributes, sortOrder, startIndex
// (1-based), and count (default 100).
// If successful, ListClients returns the clients and the total itemsPerPage of clients for
// all pages. If unsuccessful, ListClients returns the error.
func (a *API) ListClients(filter string, sortBy string, sortOrder SortOrder, startIndex int, itemsPerPage int) ([]Client, Page, error) {
	u := urlWithPath(*a.TargetURL, ClientsEndpoint)
	query := url.Values{}
	if filter != "" {
		query.Set("filter", filter)
	}
	if sortBy != "" {
		query.Set("sortBy", sortBy)
	}
	if sortOrder != "" {
		query.Set("sortOrder", string(sortOrder))
	}
	if startIndex == 0 {
		startIndex = 1
	}
	query.Set("startIndex", strconv.Itoa(startIndex))
	if itemsPerPage == 0 {
		itemsPerPage = 100
	}
	query.Set("count", strconv.Itoa(itemsPerPage))
	u.RawQuery = query.Encode()

	clients := &paginatedClientList{}
	err := a.doJSON(http.MethodGet, &u, nil, clients, true)
	if err != nil {
		return nil, Page{}, err
	}
	page := Page{
		StartIndex:   clients.StartIndex,
		ItemsPerPage: clients.ItemsPerPage,
		TotalResults: clients.TotalResults,
	}
	return clients.Resources, page, err
}

// ListAllClients retrieves UAA clients
func (a *API) ListAllClients(filter string, sortBy string, sortOrder SortOrder) ([]Client, error) {
	page := Page{
		StartIndex:   1,
		ItemsPerPage: 100,
	}
	var (
		results     []Client
		currentPage []Client
		err         error
	)

	for {
		currentPage, page, err = a.ListClients(filter, sortBy, sortOrder, page.StartIndex, page.ItemsPerPage)
		if err != nil {
			return nil, err
		}
		results = append(results, currentPage...)

		if (page.StartIndex + page.ItemsPerPage) > page.TotalResults {
			break
		}
		page.StartIndex = page.StartIndex + page.ItemsPerPage
	}
	return results, nil
}
