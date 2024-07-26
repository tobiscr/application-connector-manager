---
name: Test protocol
about: Document a manual testing cycle
---

# Test details

* Tested ACM Module Version: {e.g. 1.1.5}
* Github release: {Link to GH release}
* Test execution date: {2025-07-29}
* Tester: {@your-Github-user}

# Test protocol

The detailed test flow is explained in our internal documentation.

## Pre-requisites

- [ ] Configure System Landscape and Formations in BTP Cockpit
- [ ] Access to Back Office and Storefront possible

## Test flow:

### 1. Register events in Back office

 - [ ] Passed
 - [ ] Failed
 
     Link to Github issue: {GH link}

### 2. Prepare serverless function in Kyma

 - [ ] Passed
 - [ ] Failed

     Link to Github issue: {GH link}

### 3. Trigger an event in EC

 - [ ] Passed
 - [ ] Failed

     Link to Github issue: {GH link}

### 3. Verify event was delivered and API called

 - [ ] Passed
 - [ ] Failed

     Link to Github issue: {GH link}

## Cleanup

- [ ] In the Backoffice, delete the created Destination Target and Customer.
- [ ] In BTP delete the System and the Formation.
