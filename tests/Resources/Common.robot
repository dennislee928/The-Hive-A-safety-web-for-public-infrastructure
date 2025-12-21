*** Settings ***
Documentation    Common test resources and keywords for ERH Safety System PoC
Library          RequestsLibrary
Library          Collections
Library          DateTime
Library          JSONLibrary
Library          String

*** Variables ***
${BASE_URL}      http://localhost:8080
${API_VERSION}   v1
${API_BASE}      ${BASE_URL}/api/${API_VERSION}

# Zone IDs
${ZONE_Z1}       Z1
${ZONE_Z2}       Z2
${ZONE_Z3}       Z3
${ZONE_Z4}       Z4

# Decision States
${STATE_D0}      D0
${STATE_D1}      D1
${STATE_D2}      D2
${STATE_D3}      D3
${STATE_D4}      D4
${STATE_D5}      D5
${STATE_D6}      D6

# HTTP Status Codes
${HTTP_OK}              200
${HTTP_CREATED}         201
${HTTP_BAD_REQUEST}     400
${HTTP_UNAUTHORIZED}    401
${HTTP_NOT_FOUND}       404
${HTTP_INTERNAL_ERROR}  500

*** Keywords ***
Get Health Status
    [Documentation]    Check if the server is healthy
    ${response}=       GET    ${BASE_URL}/health    expected_status=any
    [Return]           ${response}

Create Test Signal
    [Documentation]    Create a test signal with given parameters
    [Arguments]        ${signal_type}    ${zone_id}    ${content}=Test signal
    ${timestamp}=      Get Current Date    result_format=%Y-%m-%dT%H:%M:%S%z
    ${signal}=         Create Dictionary
    ...                type=${signal_type}
    ...                zone_id=${zone_id}
    ...                content=${content}
    ...                timestamp=${timestamp}
    [Return]           ${signal}

Submit Crowd Report
    [Documentation]    Submit a crowd report
    [Arguments]        ${zone_id}    ${content}    ${device_id}=test_device_001
    ${headers}=        Create Dictionary    Content-Type=application/json
    ${body}=           Create Dictionary
    ...                zone_id=${zone_id}
    ...                content=${content}
    ...                device_id=${device_id}
    ${response}=       POST    ${API_BASE}/reports    json=${body}    headers=${headers}    expected_status=any
    [Return]           ${response}

Get Zone State
    [Documentation]    Get the latest decision state for a zone
    [Arguments]        ${zone_id}
    ${response}=       GET    ${API_BASE}/operator/zones/${zone_id}/state    expected_status=any
    [Return]           ${response}

Create PreAlert
    [Documentation]    Create a D0 pre-alert for a zone
    [Arguments]        ${zone_id}    ${reason}=Test pre-alert
    ${headers}=        Create Dictionary    Content-Type=application/json
    ${body}=           Create Dictionary    reason=${reason}
    ${response}=       POST    ${API_BASE}/operator/decisions/${zone_id}/d0    json=${body}    headers=${headers}    expected_status=any
    [Return]           ${response}

Verify Response Status
    [Documentation]    Verify HTTP response status
    [Arguments]        ${response}    ${expected_status}
    ${status}=         Convert To Integer    ${response.status_code}
    Should Be Equal    ${status}    ${expected_status}    Response status should be ${expected_status}

Verify Response JSON
    [Documentation]    Verify response contains valid JSON with expected fields
    [Arguments]        ${response}    ${expected_fields}=${EMPTY}
    ${json}=           Set Variable    ${response.json()}
    Should Not Be Empty    ${json}
    IF    '${expected_fields}' != '${EMPTY}'
        ${fields}=     Split String    ${expected_fields}    ,
        FOR    ${field}    IN    @{fields}
            Dictionary Should Contain Key    ${json}    ${field}
        END
    END
    [Return]           ${json}

Calculate Time Difference
    [Documentation]    Calculate time difference in seconds between two timestamps
    [Arguments]        ${start_time}    ${end_time}
    ${start}=          Parse Date    ${start_time}    result_format=%Y-%m-%dT%H:%M:%S
    ${end}=            Parse Date    ${end_time}    result_format=%Y-%m-%dT%H:%M:%S
    ${diff}=           Subtract Date From Date    ${end}    ${start}    result_format=number
    ${seconds}=        Evaluate    ${diff} * 86400
    [Return]           ${seconds}

Wait For State Transition
    [Documentation]    Wait for a zone to transition to a specific state
    [Arguments]        ${zone_id}    ${expected_state}    ${timeout}=30
    ${end_time}=       Evaluate    time.time() + ${timeout}    modules=time
    ${found}=          Set Variable    ${False}
    WHILE    time.time() < ${end_time} and not ${found}
        ${response}=   Get Zone State    ${zone_id}
        ${status}=     Set Variable    ${response.status_code}
        Run Keyword If    ${status} == 200
        ...    Check State Transition    ${response}    ${expected_state}    ${found}
    END
    Should Be True    ${found}    Zone ${zone_id} did not transition to ${expected_state} within ${timeout} seconds

Check State Transition
    [Documentation]    Check if zone is in expected state
    [Arguments]        ${response}    ${expected_state}    ${found}
    ${json}=           Set Variable    ${response.json()}
    ${current_state}=  Get From Dictionary    ${json}    current_state
    ${found}=          Set Variable If    '${current_state}' == '${expected_state}'    ${True}    ${False}
    [Return]           ${found}

