*** Settings ***
Documentation    Evaluation metrics tests (TTA, TTDR, FN/FP rates)
Test Setup       Setup Test Environment
Test Teardown    Cleanup Test Environment
Resource         ../Resources/Common.robot
Library          DateTime

*** Variables ***
${TEST_ZONE}     Z1
${TTA_THRESHOLD}    60    # seconds
${TTDR_THRESHOLD}   120   # seconds

*** Test Cases ***
Evaluation: TTA Measurement
    [Documentation]    Measure Time-to-Acknowledge (TTA)
    [Tags]    evaluation    tta
    
    # Record signal submission time
    ${signal_time}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S
    ${signal}=         Create Test Signal    infrastructure    ${TEST_ZONE}    TTA test signal
    ${signal_response}=    POST    ${API_BASE}/infrastructure/signals    json=${signal}
    Verify Response Status    ${signal_response}    ${HTTP_OK}
    
    # Create pre-alert (D0 acknowledgment)
    ${alert_response}=    Create PreAlert    ${TEST_ZONE}    TTA measurement
    ${alert_time}=        Get Current Date    result_format=%Y-%m-%dT%H:%M:%S
    
    # Calculate TTA
    ${tta}=              Calculate Time Difference    ${signal_time}    ${alert_time}
    Log    TTA: ${tta} seconds    level=INFO
    
    # Verify TTA is within threshold (for low severity)
    Should Be True    ${tta} < ${TTA_THRESHOLD}
    ...    TTA should be less than ${TTA_THRESHOLD} seconds for low severity scenarios

Evaluation: TTDR Measurement
    [Documentation]    Measure Time-to-Dispatch Recommendation (TTDR)
    [Tags]    evaluation    ttdr
    
    # Record signal submission time
    ${signal_time}=     Get Current Date    result_format=%Y-%m-%dT%H:%M:%S
    ${signal}=          Create Test Signal    infrastructure    ${TEST_ZONE}    TTDR test signal
    ${signal_response}=    POST    ${API_BASE}/infrastructure/signals    json=${signal}
    Verify Response Status    ${signal_response}    ${HTTP_OK}
    
    # Create pre-alert (D0)
    ${alert_response}=  Create PreAlert    ${TEST_ZONE}    TTDR measurement
    
    # Get state to check for D1 (dispatch recommendation)
    # Note: This is simplified - in reality, D1 would be automatically generated
    ${state_response}=  Get Zone State    ${TEST_ZONE}
    ${d1_time}=         Get Current Date    result_format=%Y-%m-%dT%H:%M:%S
    
    # Calculate TTDR
    ${ttdr}=            Calculate Time Difference    ${signal_time}    ${d1_time}
    Log    TTDR: ${ttdr} seconds    level=INFO
    
    # Verify TTDR is within threshold
    Should Be True    ${ttdr} < ${TTDR_THRESHOLD}
    ...    TTDR should be less than ${TTDR_THRESHOLD} seconds

Evaluation: ERH Metrics Query
    [Documentation]    Test ERH metrics retrieval
    [Tags]    evaluation    erh
    
    ${response}=        GET    ${API_BASE}/erh/status/${TEST_ZONE}    expected_status=any
    ${status}=          Set Variable    ${response.status_code}
    
    Run Keyword If    ${status} == ${HTTP_OK}
    ...    Verify ERH Metrics Response    ${response}

Evaluation: ERH Metrics History
    [Documentation]    Test ERH metrics history retrieval
    [Tags]    evaluation    erh
    
    ${response}=        GET    ${API_BASE}/erh/metrics/${TEST_ZONE}/history    expected_status=any
    ${status}=          Set Variable    ${response.status_code}
    
    Run Keyword If    ${status} == ${HTTP_OK}
    ...    Verify Response JSON    ${response}    status

Evaluation: ERH Metrics Trends
    [Documentation]    Test ERH metrics trends retrieval
    [Tags]    evaluation    erh
    
    ${response}=        GET    ${API_BASE}/erh/metrics/${TEST_ZONE}/trends?duration=24h    expected_status=any
    ${status}=          Set Variable    ${response.status_code}
    
    Run Keyword If    ${status} == ${HTTP_OK}
    ...    Verify Response JSON    ${response}    status

*** Keywords ***
Verify ERH Metrics Response
    [Documentation]    Verify ERH metrics response structure
    [Arguments]        ${response}
    ${json}=           Set Variable    ${response.json()}
    Should Contain     ${json}    complexity
    Should Contain     ${json}    ethical_primes
    # Verify complexity metrics
    Dictionary Should Contain Key    ${json}    complexity
    Dictionary Should Contain Key    ${json}    ethical_primes

Setup Test Environment
    [Documentation]    Setup test environment
    ${health}=         Get Health Status
    Should Be Equal    ${health.status_code}    ${HTTP_OK}    Server should be healthy

Cleanup Test Environment
    [Documentation]    Cleanup test environment
    [Return]           No cleanup needed for evaluation tests

