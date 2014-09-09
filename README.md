nstop or (!stop)
=====

[![Build Status](https://api.shippable.com/projects/540efe883479c5ea8f9f2b21/badge?branchName=master)](https://app.shippable.com/projects/540efe883479c5ea8f9f2b21/builds/latest)

docker container process manager

Why build another supervisor?
=============================

Mostly to solve the problem's we have been facing in managing our docker-ized applications like

- restarting the apps upon crash  
- not being able to force restart an application 
- auto-restarting application upon file changes 
- not being able to handle logs 
- other solutions like node-supervisor, supervisord, forever etc had one or more shortcomings

How nstop is different?
=============================

- built from ground up to address only containerized app specific issues
- one stop solution for all container specific needs like process manager, file watcher, log rotate
- different from phusion/baseimage because baseimage solves the problem of running multiple processes inside a container and basically make it
easier to treat a container like a VM whereas nstop sticks to docker philosophy of having one application per container

Want to contribue?
============================

Just open an issue and lets get going :)  

