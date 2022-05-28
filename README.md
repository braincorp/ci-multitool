# ci-multitool

go command line (and library) for interacting with apis and data in CI.

### parse pulumi output and post to PR

```
ci-multitool pulumi jsonoutput pulumi/jsonoutput/testdata/preview-changes.json -d 'gh-pr-trailer,stdout' --key pulumi-preview --repo alexgartner-bc/test --pr 3
```

![image](https://user-images.githubusercontent.com/74934191/170845809-1d2fe713-4f7f-4b57-a1e5-df3a19298fab.png)


### stdin to gihub pr

```
echo asdf | ci-multitool github pr-trailer --repo alexgartner-bc/test --pr 3 --summary "my summary" -
```

### stdin to gihub coment

```
echo asdf | ci-multitool github comment --repo alexgartner-bc/test --pr 3 -
```
