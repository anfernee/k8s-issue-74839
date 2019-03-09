FROM scratch

ADD k8s-issue-74839 /k8s-issue-74839

ENTRYPOINT ["/k8s-issue-74839"]

