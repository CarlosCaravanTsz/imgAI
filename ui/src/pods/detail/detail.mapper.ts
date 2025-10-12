import * as vm from "./detail.vm";
import * as api from "./detail.api-model";

export const mapMemberDetailFromApiToVm = (member: api.MemberDetailEntityApi): vm.MemberDetailEntity => ({
  id: member.id.toString(),
  login: member.login,
  name: member.name,
  company: member.company,
  bio: member.bio
})

