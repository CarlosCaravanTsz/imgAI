import { mapMemberDetailFromApiToVm } from './detail.mapper'
import { MemberDetailEntity } from './detail.vm';
import { getMemberDetail as getMemberDetailApi } from './detail.api'

export const getMemberCollection = (id: string): Promise<MemberDetailEntity> => {
  return new Promise<MemberDetailEntity>(
    (resolve) => {
      getMemberDetailApi(id)
        .then((res) => {
          resolve(mapMemberDetailFromApiToVm(res))
        } )
    }
  )
}