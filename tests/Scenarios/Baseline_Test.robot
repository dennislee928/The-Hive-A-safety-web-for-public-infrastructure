*** Settings ***
Documentation    Baseline tests without crowd signals
Test Setup       Setup Test Environment
Test Teardown    Cleanup Test Environment
Resource         ../Resources/Common.robot

*** Variables ***
${TEST_ZONE}     Z1

*** Test Cases ***
Baseline: Infrastructure Signal Only
    [Documentation]    Test system response with only infrastructure signals (no crowd signals)
    [Tags]    baseline    infrastructure
    
    # Submit infrastructure signal
    ${signal}=         Create Test Signal    infrastructure    ${TEST_ZONE}    Smoke detected
    ${response}=       POST    ${API_BASE}/infrastructure/signals    json=${signal}    expected_status=any
    Verify Response Status    ${response}    ${HTTP_OK}
    
    # Verify signal was received
    ${json}=           Verify Response JSON    ${response}    status
    
    # Check if system can make decisions without crowd signals
    ${state_response}=    Get Zone State    ${TEST_ZONE}
    ${status}=         Set Variable    ${state_response.status_code}
    Should Be True    ${status} == ${HTTP_OK} or ${status} == ${HTTP_NOT_FOUND}

Baseline: Staff Signal Only
    [Documentation]    Test system response with only staff signals
    [Tags]    baseline    staff
    
    # Submit staff report
    ${report}=         Create Dictionary
    ...                zone_id=${TEST_ZONE}
    ...                content=Suspicious activity observed
    ...                staff_id=staff_001
    ${response}=       POST    ${API_BASE}/staff/reports    json=${report}    expected_status=any
    Verify Response Status    ${response}    ${HTTP_OK}
    
    # Verify response
    ${json}=           Verify Response JSON    ${response}    status

Baseline: Emergency Call Only
    [Documentation]    Test system response with only emergency calls
    [Tags]    baseline    emergency
    
    # Submit emergency call
    ${call}=           Create Dictionary
    ...                zone_id=${TEST_ZONE}
    ...                content=Medical emergency
    ...                caller_id=caller_001
    ...                priority=high
    ${response}=       POST    ${API_BASE}/emergency/calls    json=${call}    expected_status=any
    Verify Response Status    ${response}    ${HTTP_OK}
    
    # Verify response
    ${json}=           Verify Response JSON    ${response}    status

Baseline: Decision Without Crowd
    [Documentation]    Test decision making without crowd signals (baseline performance)
    [Tags]    baseline    decision
    
    # Create pre-alert (D0)
    ${response}=       Create PreAlert    ${TEST_ZONE}    Baseline test
    ${status}=         Set Variable    ${response.status_code}
    Run Keyword If    ${status} == ${HTTP_OK}
    ...    Verify Response JSON    ${response}    status
    
    # Get zone state
    ${state_response}=    Get Zone State    ${TEST_ZONE}
    ${status}=         Set Variable    ${state_response.status_code}
    Should Be True    ${status} == ${HTTP_OK} or ${status} == ${HTTP_NOT_FOUND}

*** Keywords ***
Setup Test Environment
    [Documentation]    Setup test environment
    ${health}=         Get Health Status
    Should Be Equal    ${health.status_code}    ${HTTP_OK}    Server should be healthy

Cleanup Test Environment
    [Documentation]    Cleanup test environment
    [Return]           No cleanup needed for baseline tests

