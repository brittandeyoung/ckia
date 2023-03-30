# Available Checks

| Id | Provider | Check Category |  Name | Rule Description|
|--|----------|------------|--------------------------------------|------------------------------------------------------------------------|
| ckia:aws:cost:IdleDBInstance | AWS      | Cost Optimization | RDS Idle DB Instances |  Any RDS DB instance that has not had a connection in the last 7 days is considered idle. |    
| ckia:aws:cost:UnassociatedElasticIPAddresses | AWS      | Cost Optimization | Unassociated Elastic IP Addresses |  An allocated Elastic IP address (EIP) is not associated with a running Amazon EC2 instance. |         
| ckia:aws:cost:UnderutilizedEBSVolume | AWS      | Cost Optimization | Underutilized Amazon EBS Volumes |  A volume is unattached or had less than 1 IOPS per day for the past 7 days. |         
| ckia:aws:security:RootAccountMissingMFA | AWS      | Security | MFA on Root Account |  MFA is not enabled on the root account. | 