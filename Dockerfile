FROM scratch
COPY scan /
ENTRYPOINT ["/scan"]
