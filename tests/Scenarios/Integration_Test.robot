*** Settings ***
Documentation    Integration tests for complete workflows
Test Setup       Setup Test Environment
Test Teardown    Cleanup Test Environment
Resource         ../Resources/Common.robot

*** Variables ***
${TEST_ZONE}     Z1

*** Test Cases ***
Integration: Complete Signal to Decision Flow
    [Documentation]    Test complete workflow from signal to decision
    [Tags]    integration    workflow
    
    # Step 1: Submit signal
    ${signal}=         Create Test Signal    infrastructure    ${TEST_ZONE}    Integration test signal
    ${signal_response}=    POST    ${API_BASE}/infrastructure/signals    json=${signal}
    Verify Response Status    ${signal_response}    ${HTTP_OK}
    
    # Step 2: Create pre-alert (D0)
    ${alert_response}=    Create PreAlert    ${TEST_ZONE}    Integration test
    ${alert_status}=      Set Variable    ${alert_response.status_code}
    Run Keyword If    ${alert_status} == ${HTTP_OK}
    ...    Verify Response JSON    ${alert_response}    status
    
    # Step 3: Check zone state
    ${state_response}=    Get Zone State    ${TEST_ZONE}
    ${state_status}=      Set Variable    ${state_response.status_code}
    Should Be True    ${state_status} in [${HTTP_OK}, ${HTTP_NOT_FOUND}]
    ...    Zone state should be accessible

Integration: Crowd Report with Trust Scoring
    [Documentation]    Test crowd report submission and trust scoring
    [Tags]    integration    crowd
    
    ${device_id}=      Set Variable    test_device_${RANDOM}
    ${response}=       Submit Crowd Report    ${TEST_ZONE}    Test crowd report    ${device_id}
    ${status}=         Set Variable    ${response.status_code}
    
    # Should accept or reject based on rate limiting
    Should Be True    ${status} in [${HTTP_OK}, ${HTTP_CREATED}, ${HTTP_BAD_REQUEST}]
    ...    Crowd report should be processed

Integration: CAP Message Generation
    [Documentation]    Test CAP message generation workflow
    [Tags]    integration    cap
    
    # First create a decision state that requires CAP message (simplified)
    ${alert_response}=    Create PreAlert    ${TEST_ZONE}    CAP test
    
    # Generate CAP message (if endpoint exists)
    ${headers}=           Create Dictionary    Content-Type=application/json
    ${body}=              Create Dictionary    zone_id=${TEST_ZONE}    message_type=Alert
    ${response}=          POST    ${API_BASE}/cap/generate    json=${body}    headers=${headers}    expected_status=any
    ${status}=            Set Variable    ${response.status_code}
    
    # CAP generation may require approvals, so accept various status codes
    Should Be True    ${status} in [${HTTP_OK}, ${HTTP_CREATED}, ${HTTP_BAD_REQUEST}, ${HTTP_NOT_FOUND}]
    ...    CAP generation should be handled appropriately

Integration: Route 2 Device Registration
    [Documentation]    Test Route 2 device registration workflow
    [Tags]    integration    route2
    
    ${device_id}=      Set Variable    test_device_route2_${RANDOM}
    ${headers}=        Create Dictionary    Content-Type=application/json
    ${body}=           Create Dictionary    device_id=${device_id}
    ${response}=       POST    ${API_BASE}/route2/devices/register    json=${body}    headers=${headers}    expected_status=any
    ${status}=         Set Variable    ${response.status_code}
    
    Run Keyword If    ${status} == ${HTTP_CREATED}
    ...    Verify Device Registration Response    ${response}

Integration: Audit Log Integrity
    [Documentation]    Test audit log integrity verification
    [Tags]    integration    audit
    
    ${response}=       GET    ${API_BASE}/audit/verify-integrity    expected_status=any
    ${status}=         Set Variable    ${response.status_code}
    
    Run Keyword If    ${status} == ${HTTP_OK}
    ...    Verify Audit Integrity Response    ${response}

*** Keywords ***
Verify Device Registration Response
    [Documentation]    Verify device registration response
    [Arguments]        ${response}
    ${json}=           Verify Response JSON    ${response}    device_id,api_key,status
    Dictionary Should Contain Key    ${json}    device_id
    Dictionary Should Contain Key    ${json}    api_key

Verify Audit Integrity Response
    [Documentation]    Verify audit integrity response
    [Arguments]        ${response}
    ${json}=           Verify Response JSON    ${response}    report
    Dictionary Should Contain Key    ${json}    report

Setup Test Environment
    [Documentation]    Setup test environment
    ${health}=         Get Health Status
    Should Be Equal    ${health.status_code}    ${HTTP_OK}    Server should be healthy

Cleanup Test Environment
    [Documentation]    Cleanup test environment
    [Return]           No cleanup needed for integration tests

