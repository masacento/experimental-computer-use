# Expelimental Computer Use

- computer use by [RobotGO](https://github.com/go-vgo/robotgo)
- browser use by [Rod](https://github.com/go-rod/rod)


## Usage

```
export ANTHROPIC_API_KEY="sk-ant-..."
cd examples/browsercmd
go run .
```

- Accessibility permission is required on macOS [https://support.apple.com/guide/mac-help/mh43185/mac]


## WARNING

> Computer use is a beta feature. Please be aware that computer use poses unique risks that are distinct from standard API features or chat interfaces. These risks are heightened when using computer use to interact with the internet. To minimize risks, consider taking precautions such as:
>
>  - Use a dedicated virtual machine or container with minimal privileges to prevent direct system attacks or accidents.
>  - Avoid giving the model access to sensitive data, such as account login information, to prevent information theft.
>  - Limit internet access to an allowlist of domains to reduce exposure to malicious content.
>  - Ask a human to confirm decisions that may result in meaningful real-world consequences as well as any tasks requiring affirmative consent, such as accepting cookies, executing financial transactions, or agreeing to terms of service.
> In some circumstances, Claude will follow commands found in content even if it conflicts with the user’s instructions. For example, Claude instructions on webpages or contained in images may override instructions or cause Claude to make mistakes. We suggest taking precautions to isolate Claude from sensitive data and actions to avoid risks related to prompt injection.
>
> Finally, please inform end users of relevant risks and obtain their consent prior to enabling computer use in your own products.
> https://docs.anthropic.com/en/docs/build-with-claude/computer-use


## LICENSE
- [JSON License](https://json.org/license.html)
> The Software shall be used for Good, not Evil.