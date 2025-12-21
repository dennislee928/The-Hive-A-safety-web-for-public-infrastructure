*** Settings ***
Documentation    Performance and load tests
Test Setup       Setup Test Environment
Test Teardown    Cleanup Test Environment
Resource         ../Resources/Common.robot
Library          DateTime

*** Variables ***
${TEST_ZONE}     Z1
${LOAD_COUNT}    100
${CONCURRENT_USERS}    10

*** Test Cases ***
Performance: Single Signal Response Time
    [Documentation]    Measure response time for single signal submission
    [Tags]    performance    signal
    
    ${start_time}=     Get Current Date    result_format=timestamp
    ${signal}=         Create Test Signal    infrastructure    ${TEST_ZONE}    Performance test
    ${response}=       POST    ${API_BASE}/infrastructure/signals    json=${signal}
    ${end_time}=       Get Current Date    result_format=timestamp
    
    Verify Response Status    ${response}    ${HTTP_OK}
    
    ${response_time}=    Calculate Time Difference    ${start_time}    ${end_time}
    Should Be True    ${response_time} < 1.0    Response time should be less than 1 second
    
    Log    Response time: ${response_time} seconds    level=INFO

Performance: Batch Signal Processing
    [Documentation]    Test system performance with batch signal processing
    [Tags]    performance    batch
    
    ${start_time}=     Get Current Date    result_format=timestamp
    FOR    ${i}    IN RANGE    ${LOAD_COUNT}
        ${signal}=     Create Test Signal    infrastructure    ${TEST_ZONE}    Batch test ${i}
        ${response}=   POST    ${API_BASE}/infrastructure/signals    json=${signal}    expected_status=any
    END
    ${end_time}=       Get Current Date    result_format=timestamp
    
    ${total_time}=     Calculate Time Difference    ${start_time}    ${end_time}
    ${avg_time}=       Evaluate    ${total_time} / ${LOAD_COUNT}
    
    Log    Total time: ${total_time} seconds    level=INFO
    Log    Average time per signal: ${avg_time} seconds    level=INFO
    Should Be True    ${avg_time} < 0.5    Average response time should be less than 0.5 seconds

Performance: Concurrent Signal Processing
    [Documentation]    Test system performance with concurrent signal processing
    [Tags]    performance    concurrent
    
    ${start_time}=     Get Current Date    result_format=timestamp
    
    # Submit multiple signals concurrently (simulated with rapid sequential requests)
    FOR    ${i}    IN RANGE    ${CONCURRENT_USERS}
        ${signal}=     Create Test Signal    infrastructure    ${TEST_ZONE}    Concurrent test ${i}
        POST    ${API_BASE}/infrastructure/signals    json=${signal}    expected_status=any    timeout=5
    END
    
    ${end_time}=       Get Current Date    result_format=timestamp
    ${total_time}=     Calculate Time Difference    ${start_time}    ${end_time}
    
    Log    Total time for ${CONCURRENT_USERS} concurrent requests: ${total_time} seconds    level=INFO
    Should Be True    ${total_time} < 10.0    Concurrent requests should complete within 10 seconds

Performance: State Query Response Time
    [Documentation]    Measure response time for zone state queries
    [Tags]    performance    query
    
    ${start_time}=     Get Current Date    result_format=timestamp
    ${response}=       Get Zone State    ${TEST_ZONE}
    ${end_time}=       Get Current Date    result_format=timestamp
    
    ${response_time}=    Calculate Time Difference    ${start_time}    ${end_time}
    Should Be True    ${response_time} < 0.5    State query should be fast
    
    Log    State query response time: ${response_time} seconds    level=INFO

*** Keywords ***
Setup Test Environment
    [Documentation]    Setup test environment
    ${health}=         Get Health Status
    Should Be Equal    ${health.status_code}    ${HTTP_OK}    Server should be healthy

Cleanup Test Environment
    [Documentation]    Cleanup test environment
    [Return]           No cleanup needed for performance tests

