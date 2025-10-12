import { MemberDetailEntityApi } from "./detail.api-model";

export const getMemberDetail = (id: string): Promise<MemberDetailEntityApi> => {

  return fetch(`https://api.github.com/orgs/lemoncode/members/${id}`)
    .then(member => member.json())
}